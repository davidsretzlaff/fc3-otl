using Dapper;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Data.Sqlite;
using OpenTelemetry.Trace;
using System.Diagnostics;
using System.Text.Json;

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

        // Serialize result to JSON before adding to trace
        string serializedResult = JsonSerializer.Serialize(result);
        span.SetAttribute("response", serializedResult);

        return Ok(result);
    }
}