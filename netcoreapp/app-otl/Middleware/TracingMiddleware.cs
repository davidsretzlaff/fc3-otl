using System.Text;
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

        public TracingMiddleware(RequestDelegate next, Tracer tracer)
        {
            _next = next;
            _tracer = tracer;
        }

        // Método principal do middleware que é chamado para cada requisição
        public async Task InvokeAsync(HttpContext context)
        {
            // Verifica se já existe um span ativo para esta requisição
            var existingSpan = Tracer.CurrentSpan;
            if (existingSpan != null && existingSpan.IsRecording)
            {
                // Se já existe um span ativo, apenas adiciona os atributos
                var request = await FormatRequest(context.Request);
                existingSpan.SetAttribute("http.request", request);

                var originalBodyStream = context.Response.Body;
                using var responseBody = new MemoryStream();
                context.Response.Body = responseBody;

                try
                {
                    await _next(context);

                    var response = await FormatResponse(context.Response);
                    existingSpan.SetAttribute("http.response", response);
                    existingSpan.SetAttribute("http.status_code", context.Response.StatusCode.ToString());

                    await responseBody.CopyToAsync(originalBodyStream);
                }
                catch (Exception ex)
                {
                    existingSpan.SetAttribute("error", true);
                    existingSpan.SetAttribute("error.message", ex.Message);
                    throw;
                }
            }
            else
            {
                // Se não existe span ativo, cria um novo
                using var span = _tracer.StartActiveSpan($"{context.Request.Method} {context.Request.Path}");

                var request = await FormatRequest(context.Request);
                span.SetAttribute("http.request", request);

                var originalBodyStream = context.Response.Body;
                using var responseBody = new MemoryStream();
                context.Response.Body = responseBody;

                try
                {
                    await _next(context);

                    var response = await FormatResponse(context.Response);
                    span.SetAttribute("http.response", response);
                    span.SetAttribute("http.status_code", context.Response.StatusCode.ToString());

                    await responseBody.CopyToAsync(originalBodyStream);
                }
                catch (Exception ex)
                {
                    span.SetAttribute("error", true);
                    span.SetAttribute("error.message", ex.Message);
                    throw;
                }
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
            var headers = string.Join(", ", request.Headers.Select(h => $"{h.Key}: {h.Value}"));
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
            var headers = string.Join(", ", response.Headers.Select(h => $"{h.Key}: {h.Value}"));
            // Retorna uma string formatada com todas as informações da resposta
            return $"Status: {response.StatusCode}, Headers: {headers}, Body: {body}";
        }
    }
} 