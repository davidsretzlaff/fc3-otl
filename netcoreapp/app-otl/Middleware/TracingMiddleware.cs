using System.Text;
using System.Text.RegularExpressions;
using Microsoft.Extensions.Primitives;
using OpenTelemetry.Trace;

namespace app_otl.Middleware
{
    public class TracingMiddleware
    {
        private readonly RequestDelegate _next;
        private readonly Tracer _tracer;

        // Regex para identificar cartão de crédito
        private static readonly Regex CreditCardRegex = new(@"\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b");

        public TracingMiddleware(RequestDelegate next, Tracer tracer)
        {
            _next = next;
            _tracer = tracer;
        }

        public async Task InvokeAsync(HttpContext context)
        {
            // Cria um span para a requisição
            using var span = _tracer.StartActiveSpan($"{context.Request.Method} {context.Request.Path}");

            // Captura e mascara informações do request
            var request = await FormatRequest(context.Request);
            request = MaskSensitiveData(request);
            span.SetAttribute("http.request", request);

            // Captura o response original
            var originalBodyStream = context.Response.Body;
            using var responseBody = new MemoryStream();
            context.Response.Body = responseBody;

            try
            {
                await _next(context);

                // Captura e mascara informações do response
                var response = await FormatResponse(context.Response);
                response = MaskSensitiveData(response);
                span.SetAttribute("http.response", response);
                span.SetAttribute("http.status_code", context.Response.StatusCode.ToString());

                // Copia o response para o stream original
                await responseBody.CopyToAsync(originalBodyStream);
            }
            catch (Exception ex)
            {
                span.SetAttribute("error", true);
                span.SetAttribute("error.message", ex.Message);
                throw;
            }
        }

        private async Task<string> FormatRequest(HttpRequest request)
        {
            request.EnableBuffering();
            var body = await new StreamReader(request.Body).ReadToEndAsync();
            request.Body.Position = 0;

            var headers = string.Join(", ", request.Headers.Select(h => $"{h.Key}: {MaskSensitiveData(h.Value)}"));
            return $"Method: {request.Method}, Path: {request.Path}, Headers: {headers}, Query: {request.QueryString}, Body: {body}";
        }

        private async Task<string> FormatResponse(HttpResponse response)
        {
            response.Body.Seek(0, SeekOrigin.Begin);
            var body = await new StreamReader(response.Body).ReadToEndAsync();
            response.Body.Seek(0, SeekOrigin.Begin);

            var headers = string.Join(", ", response.Headers.Select(h => $"{h.Key}: {MaskSensitiveData(h.Value)}"));
            return $"Status: {response.StatusCode}, Headers: {headers}, Body: {body}";
        }

        private string MaskSensitiveData(string input)
        {
            if (string.IsNullOrEmpty(input))
                return input;

            // Mascara cartão de crédito
            return CreditCardRegex.Replace(input, match =>
            {
                var card = match.Value.Replace(" ", "").Replace("-", "");
                return $"{card[..4]}-XXXX-XXXX-{card[^4..]}";
            });
        }

        private string MaskSensitiveData(StringValues values)
        {
            return string.Join(", ", values.Select(v => MaskSensitiveData(v)));
        }
    }
} 