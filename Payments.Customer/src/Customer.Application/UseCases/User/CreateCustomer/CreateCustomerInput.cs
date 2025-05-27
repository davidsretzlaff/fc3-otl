using Customer.Application.UseCases.User.Common;
using MediatR;

namespace Customer.Application.UseCases.User.CreateUser;

public record CreateCustomerInput : IRequest<CustomerOutput>
{
    public string Name { get; init; }
    public string Email { get; init; }

    public CreateCustomerInput(string name, string email)
    {
        Name = name;
        Email = email;
    }
}