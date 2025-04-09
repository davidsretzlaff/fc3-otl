using System.Text;
using System.Text.RegularExpressions;
using Microsoft.Extensions.Primitives;
using OpenTelemetry.Trace;

namespace app_otl.Middleware
{
    public class TracingMiddleware
    {
        // _next: Representa o próximo middleware na pipeline
        // _tracer: Instância do OpenTelemetry Tracer para criar spans
        private readonly RequestDelegate _next;
        private readonly Tracer _tracer;

        // Regex para identificar números de cartão de crédito no formato:
        // XXXX-XXXX-XXXX-XXXX ou XXXX XXXX XXXX XXXX
        private static readonly Regex CreditCardRegex = new(@"\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b");

        public TracingMiddleware(RequestDelegate next, Tracer tracer)
        {
            _next = next;
            _tracer = tracer;
        }

        // Método principal do middleware que é chamado para cada requisição
        public async Task InvokeAsync(HttpContext context)
        {
            // 1. Cria um novo span para rastrear esta requisição
            // O nome do span é composto pelo método HTTP e caminho (ex: "GET /MyController")
            using var span = _tracer.StartActiveSpan($"{context.Request.Method} {context.Request.Path}");

            // 2. Captura e formata as informações da requisição
            // Inclui método, path, headers, query string e body
            var request = await FormatRequest(context.Request);
            // 3. Mascara dados sensíveis no request (como cartões de crédito)
            request = MaskSensitiveData(request);
            // 4. Adiciona o request formatado como atributo no span
            span.SetAttribute("http.request", request);

            // 5. Prepara para capturar o response
            // Salva o stream original do response
            var originalBodyStream = context.Response.Body;
            // Cria um novo MemoryStream para capturar o response
            using var responseBody = new MemoryStream();
            // Substitui o stream do response pelo nosso MemoryStream
            context.Response.Body = responseBody;

            try
            {
                // 6. Chama o próximo middleware na pipeline
                // Isso permite que a requisição continue seu processamento normal
                await _next(context);

                // 7. Após o processamento, captura e formata o response
                var response = await FormatResponse(context.Response);
                // 8. Mascara dados sensíveis no response
                response = MaskSensitiveData(response);
                // 9. Adiciona o response formatado como atributo no span
                span.SetAttribute("http.response", response);
                // 10. Adiciona o status code como atributo
                span.SetAttribute("http.status_code", context.Response.StatusCode.ToString());

                // 11. Copia o response processado de volta para o stream original
                // Isso garante que o cliente receba a resposta correta
                await responseBody.CopyToAsync(originalBodyStream);
            }
            catch (Exception ex)
            {
                // 12. Se ocorrer algum erro, marca o span como erro
                span.SetAttribute("error", true);
                span.SetAttribute("error.message", ex.Message);
                throw;
            }
        }

        // Formata as informações da requisição
        private async Task<string> FormatRequest(HttpRequest request)
        {
            // Habilita o buffering para permitir múltiplas leituras do body
            request.EnableBuffering();
            // Lê o corpo da requisição
            var body = await new StreamReader(request.Body).ReadToEndAsync();
            // Reseta a posição do stream para permitir que outros middlewares leiam
            request.Body.Position = 0;

            // Formata os headers, aplicando máscara em dados sensíveis
            var headers = string.Join(", ", request.Headers.Select(h => $"{h.Key}: {MaskSensitiveData(h.Value)}"));
            // Retorna uma string formatada com todas as informações da requisição
            return $"Method: {request.Method}, Path: {request.Path}, Headers: {headers}, Query: {request.QueryString}, Body: {body}";
        }

        // Formata as informações da resposta
        private async Task<string> FormatResponse(HttpResponse response)
        {
            // Move o ponteiro do stream para o início
            response.Body.Seek(0, SeekOrigin.Begin);
            // Lê o corpo da resposta
            var body = await new StreamReader(response.Body).ReadToEndAsync();
            // Reseta a posição do stream
            response.Body.Seek(0, SeekOrigin.Begin);

            // Formata os headers, aplicando máscara em dados sensíveis
            var headers = string.Join(", ", response.Headers.Select(h => $"{h.Key}: {MaskSensitiveData(h.Value)}"));
            // Retorna uma string formatada com todas as informações da resposta
            return $"Status: {response.StatusCode}, Headers: {headers}, Body: {body}";
        }

        // Mascara dados sensíveis em uma string
        private string MaskSensitiveData(string input)
        {
            if (string.IsNullOrEmpty(input))
                return input;

            // Procura por números de cartão de crédito e os mascara
            // Exemplo: 4532-1234-5678-9012 -> 4532-XXXX-XXXX-9012
            return CreditCardRegex.Replace(input, match =>
            {
                // Remove espaços e hífens do número do cartão
                var card = match.Value.Replace(" ", "").Replace("-", "");
                // Mantém os primeiros 4 e últimos 4 dígitos, mascara o resto
                return $"{card[..4]}-XXXX-XXXX-{card[^4..]}";
            });
        }

        // Mascara dados sensíveis em um conjunto de valores (usado para headers)
        private string MaskSensitiveData(StringValues values)
        {
            return string.Join(", ", values.Select(v => MaskSensitiveData(v)));
        }
    }
} 