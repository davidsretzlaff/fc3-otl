using Microsoft.Data.Sqlite;
using Dapper;
using Customer.Domain.Domain;
using Customer.Application.Interfaces;
using Microsoft.Extensions.Logging;

namespace Customer.Infra.Data.Repositories;

public class CustomerRepository : ICustomerRepository
{
    private readonly string _connectionString = "Data Source=/app/data/customer.db";
    private readonly List<Domain.Domain.Customer> _customers = new();
    private readonly ILogger<CustomerRepository> _logger;

    public CustomerRepository(ILogger<CustomerRepository> logger)
    {
        _logger = logger;
    }

    public async Task Add(Domain.Domain.Customer user)
    {
        _customers.Add(user);
        _logger.LogInformation("[customer] Customer added to in-memory collection. Email: {Email}", user.Email);
    }

    public async Task CreateCustomer(Domain.Domain.Customer user, CancellationToken cancellationToken)
    {
        _logger.LogInformation("[customer] Attempting to create customer in database. Email: {Email}", user.Email);

        try
        {
            using var connection = new SqliteConnection(_connectionString);
            await connection.OpenAsync(cancellationToken);
            
            var sql = @"
                INSERT INTO Customers (Id, Name, Email, CreatedAt, UpdatedAt)
                VALUES (@Id, @Name, @Email, @CreatedAt, @UpdatedAt)";

            await connection.ExecuteAsync(sql, new
            {
                Id = user.Id.ToString(),
                Name = user.Name,
                Email = user.Email,
                CreatedAt = DateTime.UtcNow.ToString("yyyy-MM-dd HH:mm:ss"),
                UpdatedAt = DateTime.UtcNow.ToString("yyyy-MM-dd HH:mm:ss")
            });

            _logger.LogInformation("[customer] Customer created successfully in database. Email: {Email}", user.Email);
        }
        catch (SqliteException ex)
        {
            _logger.LogError("[customer] ERROR: SQLite database error - {ErrorMessage}. Email: {Email}", 
                ex.Message, user.Email);
            throw;
        }
        catch (Exception ex)
        {
            _logger.LogError("[customer] ERROR: Database connection failed - {ErrorMessage}. Email: {Email}", 
                ex.Message, user.Email);
            throw;
        }
    }

    public async Task<Domain.Domain.Customer> GetUserById(Guid id, CancellationToken cancellationToken)
    {
        _logger.LogInformation("[customer] Retrieving customer from database. CustomerId: {CustomerId}", id);

        try
        {
            using var connection = new SqliteConnection(_connectionString);
            await connection.OpenAsync(cancellationToken);
            
            var sql = "SELECT Id, Name, Email FROM Customers WHERE Id = @Id";
            var result = await connection.QueryFirstOrDefaultAsync(sql, new { Id = id.ToString() });
            
            if (result == null)
            {
                _logger.LogWarning("[customer] Customer not found in database. CustomerId: {CustomerId}", id);
                return null;
            }

            _logger.LogInformation("[customer] Customer retrieved successfully. CustomerId: {CustomerId}", id);
            return Domain.Domain.Customer.Create(result.Name, result.Email);
        }
        catch (SqliteException ex)
        {
            _logger.LogError("[customer] ERROR: SQLite error during retrieval - {ErrorMessage}. CustomerId: {CustomerId}", 
                ex.Message, id);
            throw;
        }
        catch (Exception ex)
        {
            _logger.LogError("[customer] ERROR: Database error during retrieval - {ErrorMessage}. CustomerId: {CustomerId}", 
                ex.Message, id);
            throw;
        }
    }

    public async Task Save()
    {
        _logger.LogInformation("[customer] Saving {CustomerCount} customers to database", _customers.Count);

        try
        {
            foreach (var customer in _customers)
            {
                await CreateCustomer(customer, CancellationToken.None);
            }
            _customers.Clear();
            
            _logger.LogInformation("[customer] All customers saved successfully");
        }
        catch (Exception ex)
        {
            _logger.LogError("[customer] ERROR during batch save: {ErrorMessage}", ex.Message);
            _customers.Clear();
            throw;
        }
    }
} 