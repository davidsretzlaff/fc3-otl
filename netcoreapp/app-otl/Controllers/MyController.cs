using Microsoft.AspNetCore.Mvc;
using Microsoft.Data.Sqlite;
using Dapper;
using System.Text.Json;
using OpenTelemetry.Trace;
using System.Diagnostics;

namespace app_otl.Controllers
{
    [ApiController]
    [Route("MyController")]
    public class MyController : ControllerBase
    {
        private readonly Tracer _tracer;

        public MyController(Tracer tracer)
        {
            _tracer = tracer;
        }

        [HttpGet]
        public async Task<IActionResult> Get()
        {
            using var span = _tracer.StartActiveSpan("GetEndpoint");
            span.SetAttribute("request", "Received request");

            using var connection = new SqliteConnection("Data Source=mydatabase.db");
            var result = await connection.QueryAsync("SELECT * FROM MyTable");

            string serializedResult = JsonSerializer.Serialize(result);
            span.SetAttribute("response", serializedResult);

            return Ok(result);
        }
    }
}