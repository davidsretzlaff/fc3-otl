using Microsoft.AspNetCore.Mvc;
using OpenTelemetry.Trace;
using System.Diagnostics;

namespace app_otl.Controllers
{
    [ApiController]
    [Route("Test")]
    public class TestController : ControllerBase
    {
        private readonly ILogger<TestController> _logger;
        private readonly Tracer _tracer;

        public TestController(ILogger<TestController> logger, TracerProvider tracerProvider)
        {
            _logger = logger;
            _tracer = tracerProvider.GetTracer("TestController");
        }

        [HttpGet]
        public async Task<IActionResult> Get()
        {
            using var span = _tracer.StartActiveSpan("GetEndpoint");
            span.SetAttribute("request", "Received request");

            _logger.LogInformation("Rota /Test acessada");

            // Simula uma operação demorada
            await Task.Delay(100);

            using var client = new HttpClient();
            var response = await client.GetAsync("https://httpbin.org/get");
            var content = await response.Content.ReadAsStringAsync();

            span.SetAttribute("response.length", content.Length);
            return Ok(content);
        }
    }
} 