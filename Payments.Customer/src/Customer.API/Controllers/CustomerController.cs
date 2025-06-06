using app_otl.ApiModels.Response;
using Customer.Application.UseCases.User.Common;
using Customer.Application.UseCases.User.CreateUser;
using Customer.API.Controllers;
using Customer.API.Middleware;
using Dapper;
using MediatR;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Data.Sqlite;

namespace app_otl.Controllers
{
    [ApiController]
    [Route("api/customer")]  
    public class CustomerController : BaseController
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

                return SuccessResponse(output, 201, "Customer created successfully");
            }
            catch (Exception ex)
            {
                // LOG ESTRUTURADO - Erro
                _logger.LogError(ex, "Failed to create customer for {CustomerEmail}", input.Email);
                
                return InternalServerErrorResponse("An error occurred while creating the customer", exception: ex);
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
                
                return SuccessResponse(customers);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Failed to retrieve customers");
                
                return InternalServerErrorResponse("An error occurred while retrieving customers", exception: ex);
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
                    return NotFoundResponse($"Customer with ID '{id}' was not found");
                }
                
                return SuccessResponse(customer);
            }
            catch (Exception ex)
            {
                _logger.LogError(ex, "Failed to retrieve customer {CustomerId}", id);
                
                return InternalServerErrorResponse("An error occurred while retrieving the customer", exception: ex);
            }
        }
    }
}