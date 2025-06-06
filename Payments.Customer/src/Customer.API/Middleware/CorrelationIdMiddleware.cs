using System.Diagnostics;
using Serilog.Context;

namespace Customer.API.Middleware
{
    public class CorrelationIdMiddleware
    {
        private readonly RequestDelegate _next;
        private const string CorrelationIdHeaderName = "X-Correlation-ID";

        public CorrelationIdMiddleware(RequestDelegate next)
        {
            _next = next;
        }

        public async Task InvokeAsync(HttpContext context)
        {
            // Extrair ou gerar correlation ID
            var correlationId = context.Request.Headers[CorrelationIdHeaderName].FirstOrDefault();
            
            if (string.IsNullOrEmpty(correlationId))
            {
                correlationId = GenerateCorrelationId("customer");
            }

            // Adicionar ao contexto
            context.Items["CorrelationId"] = correlationId;
            context.Response.Headers[CorrelationIdHeaderName] = correlationId;

            // Adicionar ao contexto do Serilog para aparecer nos logs JSON
            using (LogContext.PushProperty("correlation_id", correlationId))
            {
                // Adicionar ao Activity (OpenTelemetry)
                var activity = Activity.Current;
                if (activity != null)
                {
                    activity.SetTag("correlation.id", correlationId);
                }

                // EXECUÇÃO SILENCIOSA - sem logs de requisição
                await _next(context);
            }
        }

        private static string GenerateCorrelationId(string service)
        {
            var timestamp = DateTime.UtcNow.ToString("yyyyMMddHHmmss");
            var random = Guid.NewGuid().ToString("N")[..8];
            return $"{service}-{timestamp}-{random}";
        }
    }

    public static class CorrelationIdMiddlewareExtensions
    {
        public static IApplicationBuilder UseCorrelationId(this IApplicationBuilder builder)
        {
            return builder.UseMiddleware<CorrelationIdMiddleware>();
        }
    }

    public static class CorrelationIdHelper
    {
        public static string? GetCorrelationId(this HttpContext context)
        {
            return context.Items["CorrelationId"] as string;
        }
    }
} 