using Customer.Application.UseCases.User.Common;
using MediatR;

namespace Customer.Application.UseCases.User.CreateUser;

public record CreateUserInput : IRequest<UserOutput>
{
    public string Name { get; init; }
    public string Login { get; init; }
    public string Password { get; init; }

    public CreateUserInput(string name, string login, string password)
    {
        Name = name;
        Login = login;
        Password = password;
    }
}