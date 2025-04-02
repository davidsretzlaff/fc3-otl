using Dapper;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Data.Sqlite;
using OpenTelemetry.Trace;
using System.Diagnostics;

[ApiController]
[Route("MyController")]
public class MyController : ControllerBase
{
    private readonly Tracer _tracer;

    public MyController(TracerProvider tracerProvider)
    {
        _tracer = tracerProvider.GetTracer("MyController");
    }

    [HttpGet]
    public async Task<IActionResult> Get()
    {
        using var span = _tracer.StartActiveSpan("GetEndpoint");
        span.SetAttribute("request", "Received request");

        using var connection = new SqliteConnection("Data Source=mydatabase.db");
        var result = await connection.QueryAsync("SELECT * FROM MyTable");

        span.SetAttribute("response", "Sending response");
        return Ok(result);
    }
}