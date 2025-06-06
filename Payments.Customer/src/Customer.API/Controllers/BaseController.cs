using app_otl.ApiModels.Response;
using Microsoft.AspNetCore.Mvc;
using System.Diagnostics;

namespace Customer.API.Controllers
{
    /// <summary>
    /// Controller base que fornece métodos padronizados para respostas da API
    /// </summary>
    public class BaseController : ControllerBase
    {
        /// <summary>
        /// Retorna uma resposta de erro padronizada com correlation_id
        /// </summary>
        protected IActionResult ErrorResponse(int statusCode, string title, string detail, Exception? exception = null)
        {
            var correlationId = GetCorrelationId();
            
            var problemDetails = new ProblemDetails
            {
                Title = title,
                Status = statusCode,
                Detail = detail,
                Instance = HttpContext.Request.Path
            };

            // Adiciona correlation_id como extensão do ProblemDetails
            problemDetails.Extensions["correlation_id"] = correlationId;
            problemDetails.Extensions["traceId"] = Activity.Current?.Id ?? HttpContext.TraceIdentifier;

            // Adiciona correlation_id no header da resposta
            HttpContext.Response.Headers.Add("X-Correlation-ID", correlationId);

            return StatusCode(statusCode, problemDetails);
        }

        /// <summary>
        /// Retorna uma resposta de sucesso padronizada com correlation_id
        /// </summary>
        protected IActionResult SuccessResponse<T>(T data, int statusCode = 200, string? message = null)
        {
            var correlationId = GetCorrelationId();
            
            var response = new ApiResponse<T>(data, message)
            {
                CorrelationId = correlationId,
                TraceId = Activity.Current?.Id ?? HttpContext.TraceIdentifier
            };

            // Adiciona correlation_id no header da resposta
            HttpContext.Response.Headers.Add("X-Correlation-ID", correlationId);

            return StatusCode(statusCode, response);
        }

        /// <summary>
        /// Retorna uma resposta de erro 400 (Bad Request) padronizada
        /// </summary>
        protected IActionResult BadRequestResponse(string detail, string title = "Bad Request")
        {
            return ErrorResponse(400, title, detail);
        }

        /// <summary>
        /// Retorna uma resposta de erro 404 (Not Found) padronizada
        /// </summary>
        protected IActionResult NotFoundResponse(string detail, string title = "Not Found")
        {
            return ErrorResponse(404, title, detail);
        }

        /// <summary>
        /// Retorna uma resposta de erro 422 (Unprocessable Entity) padronizada
        /// </summary>
        protected IActionResult UnprocessableEntityResponse(string detail, string title = "Unprocessable Entity")
        {
            return ErrorResponse(422, title, detail);
        }

        /// <summary>
        /// Retorna uma resposta de erro 500 (Internal Server Error) padronizada
        /// </summary>
        protected IActionResult InternalServerErrorResponse(string detail = "An internal server error occurred", string title = "Internal Server Error", Exception? exception = null)
        {
            return ErrorResponse(500, title, detail, exception);
        }

        /// <summary>
        /// Obtém o correlation_id do contexto atual
        /// </summary>
        private string GetCorrelationId()
        {
            return HttpContext.Items["CorrelationId"]?.ToString() ?? 
                   HttpContext.Request.Headers["X-Correlation-ID"].FirstOrDefault() ?? 
                   Guid.NewGuid().ToString();
        }
    }
} 