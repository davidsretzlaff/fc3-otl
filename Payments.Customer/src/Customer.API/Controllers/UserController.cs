using app_otl.ApiModels.Response;
using Customer.Application.UseCases.User.Common;
using Customer.Application.UseCases.User.CreateUser;
using Dapper;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Data.Sqlite;
using OpenTelemetry.Trace;

namespace app_otl.Controllers
{
    [ApiController]
    [Route("api/users")]  
    public class UserController : ControllerBase
    {
        private readonly Tracer _tracer;
        private readonly ILogger<UserController> _logger;
        private readonly IMediator _mediator;

        public UserController(Tracer tracer, ILogger<UserController> logger, IMediator mediator)
        {
            _tracer = tracer;
            _logger = logger;
            _mediator = mediator;
        }

        [HttpPost]
        [ProducesResponseType(typeof(ApiResponse<UserOutput>), StatusCodes.Status201Created)]
        [ProducesResponseType(typeof(ProblemDetails), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(ProblemDetails), StatusCodes.Status422UnprocessableEntity)]
        public async Task<IActionResult> Create([FromBody] CreateUserInput input, CancellationToken cancellationToken)
        {
            var output = await _mediator.Send(input, cancellationToken);
            return CreatedAtAction(
               nameof(Create),
               new { output.Id },
               new ApiResponse<UserOutput>(output)
            );
        }

        [HttpGet]
        public async Task<IActionResult> Get()
        {
            _logger.LogInformation("Accessing My Controller");

            _logger.LogInformation("Opening connection to database");
            using var connection = new SqliteConnection("Data Source=mydatabase.db");
            var result = await connection.QueryAsync("SELECT * FROM MyTable");
            _logger.LogInformation("Closing connection to database");
            return Ok(result);
            
        }
    }
}