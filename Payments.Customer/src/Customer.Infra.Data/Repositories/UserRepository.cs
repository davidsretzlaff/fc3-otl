using Microsoft.Data.Sqlite;
using Dapper;
using Customer.Domain.Domain;
using Customer.Application.Interfaces;

namespace Customer.Infra.Data.Repositories;

public class UserRepository : IUserRepository
{
    private readonly string _connectionString = "Data Source=mydatabase.db";
    private readonly List<User> _users = new();

    public async Task Add(User user)
    {
        _users.Add(user);
    }

    public async Task CreateUser(User user, CancellationToken cancellationToken)
    {
        using var connection = new SqliteConnection(_connectionString);
        await connection.OpenAsync(cancellationToken);
        
        var sql = @"
            INSERT INTO Users (Id, Name, Login, Password)
            VALUES (@Id, @Name, @Login, @Password)";

        await connection.ExecuteAsync(sql, new
        {
            Id = user.Id.ToString(),
            Name = user.Name,
            Login = user.Login,
            Password = user.Password
        });
    }

    public async Task<User> GetUserById(Guid id, CancellationToken cancellationToken)
    {
        using var connection = new SqliteConnection(_connectionString);
        await connection.OpenAsync(cancellationToken);
        
        var sql = "SELECT Id, Name, Login, Password FROM Users WHERE Id = @Id";
        var result = await connection.QueryFirstOrDefaultAsync(sql, new { Id = id.ToString() });
        
        if (result == null)
            return null;
            
        return User.Create(result.Name, result.Login, result.Password);
    }

    public async Task Save()
    {
        // Para cada usuário na lista temporária, salvar no banco
        foreach (var user in _users)
        {
            await CreateUser(user, CancellationToken.None);
        }
        _users.Clear();
    }
} 