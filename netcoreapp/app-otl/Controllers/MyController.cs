using Microsoft.AspNetCore.Mvc;
using Microsoft.Data.Sqlite;
using Dapper;
using OpenTelemetry.Trace;

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
            using var connection = new SqliteConnection("Data Source=mydatabase.db");
            var result = await connection.QueryAsync("SELECT * FROM MyTable");
            return Ok(result);
        }
    }
}