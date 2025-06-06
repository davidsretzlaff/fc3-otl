using app_otl.ApiModels.Response;
using Customer.Application.UseCases.User.Common;
using Customer.Application.UseCases.User.CreateUser;
using Customer.API.Middleware;
using Dapper;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Data.Sqlite;

namespace app_otl.Controllers
{
    [ApiController]
    [Route("api/customer")]  
    public class CustomerController : ControllerBase
    {
        private readonly ILogger<CustomerController> _logger;
        private readonly IMediator _mediator;

        public CustomerController(ILogger<CustomerController> logger, IMediator mediator)
        {
            _logger = logger;
            _mediator = mediator;
        }

        [HttpPost]
        [ProducesResponseType(typeof(ApiResponse<CustomerOutput>), StatusCodes.Status201Created)]
        [ProducesResponseType(typeof(ProblemDetails), StatusCodes.Status400BadRequest)]
        [ProducesResponseType(typeof(ProblemDetails), StatusCodes.Status422UnprocessableEntity)]
        public async Task<IActionResult> Create([FromBody] CreateCustomerInput input, CancellationToken cancellationToken)
        {
            // LOG ESTRUTURADO - correlation_id será incluído automaticamente
            _logger.LogInformation("Starting CreateCustomer for {CustomerEmail}", input.Email);

            try
            {
                var output = await _mediator.Send(input, cancellationToken);
                
                // LOG ESTRUTURADO - Sucesso
                _logger.LogInformation("Customer created successfully with ID {CustomerId}", output.Id);

                return CreatedAtAction(nameof(Create), new { output.Id }, new ApiResponse<CustomerOutput>(output));
            }
            catch (Exception ex)
            {
                // LOG ESTRUTURADO - Erro
                _logger.LogError(ex, "Failed to create customer for {CustomerEmail}", input.Email);
                
                return StatusCode(500, new ProblemDetails
                {
                    Title = "Internal Server Error",
                    Status = 500,
                    Detail = "An error occurred while creating the customer",
                    Instance = HttpContext.Request.Path
                });
            }
        }

        [HttpGet]
        public async Task<IActionResult> Get()
        {
            try
            {
                var connectionString = "Data Source=/app/data/customer.db";
                using var connection = new SqliteConnection(connectionString);
                
                var customers = await connection.QueryAsync<dynamic>(
                    "SELECT Id, Name, Email, CreatedAt, UpdatedAt FROM Customers ORDER BY CreatedAt DESC"
                );
                
                return Ok(new ApiResponse<IEnumerable<dynamic>>(customers));
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Failed to retrieve customers");
                
                return StatusCode(500, new ProblemDetails
                {
                    Title = "Internal Server Error", 
                    Status = 500,
                    Detail = "An error occurred while retrieving customers",
                    Instance = HttpContext.Request.Path
                });
            }
        }

        [HttpGet("{id}")]
        public async Task<IActionResult> GetById(string id)
        {
            try
            {
                var connectionString = "Data Source=/app/data/customer.db";
                using var connection = new SqliteConnection(connectionString);
                
                var customer = await connection.QueryFirstOrDefaultAsync<dynamic>(
                    "SELECT Id, Name, Email, CreatedAt, UpdatedAt FROM Customers WHERE Id = @Id", 
                    new { Id = id }
                );
                
                if (customer == null)
                {
                    return NotFound(new ProblemDetails
                    {
                        Title = "Customer Not Found",
                        Status = 404,
                        Detail = $"Customer with ID '{id}' was not found",
                        Instance = HttpContext.Request.Path
                    });
                }
                
                return Ok(new ApiResponse<dynamic>(customer));
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Failed to retrieve customer {CustomerId}", id);
                
                return StatusCode(500, new ProblemDetails
                {
                    Title = "Internal Server Error",
                    Status = 500,
                    Detail = "An error occurred while retrieving the customer",
                    Instance = HttpContext.Request.Path
                });
            }
        }
    }
}