using Microsoft.Data.Sqlite;
using Dapper;
using Customer.Domain.Domain;
using Customer.Application.Interfaces;

namespace Customer.Infra.Data.Repositories;

public class CustomerRepository : ICustomerRepository
{
    private readonly string _connectionString = "Data Source=mydatabase.db";
    private readonly List<Domain.Domain.Customer> _customers = new();

    public async Task Add(Domain.Domain.Customer user)
    {
        _customers.Add(user);
    }

    public async Task CreateCustomer(Domain.Domain.Customer user, CancellationToken cancellationToken)
    {
        using var connection = new SqliteConnection(_connectionString);
        await connection.OpenAsync(cancellationToken);
        
        var sql = @"
            INSERT INTO Customer (Id, Name, Email)
            VALUES (@Id, @Name, @Email)";

        await connection.ExecuteAsync(sql, new
        {
            Id = user.Id.ToString(),
            Name = user.Name,
            Email = user.Email
        });
    }

    public async Task<Domain.Domain.Customer> GetUserById(Guid id, CancellationToken cancellationToken)
    {
        using var connection = new SqliteConnection(_connectionString);
        await connection.OpenAsync(cancellationToken);
        
        var sql = "SELECT Id, Name, Email FROM Customer WHERE Id = @Id";
        var result = await connection.QueryFirstOrDefaultAsync(sql, new { Id = id.ToString() });
        
        if (result == null)
            return null;
            
        return Domain.Domain.Customer.Create(result.Name, result.eMAIL);
    }

    public async Task Save()
    {
        // Para cada usuário na lista temporária, salvar no banco
        foreach (var customer in _customers)
        {
            await CreateCustomer(customer, CancellationToken.None);
        }
        _customers.Clear();
    }
} 