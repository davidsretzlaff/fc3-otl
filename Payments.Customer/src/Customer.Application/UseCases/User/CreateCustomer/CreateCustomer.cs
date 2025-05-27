using Customer.Application.Interfaces;
using Customer.Application.UseCases.User.Common;
using MediatR;
using DomainEntity = Customer.Domain.Domain;

namespace Customer.Application.UseCases.User.CreateUser;

public class CreateCustomer : IRequestHandler<CreateCustomerInput, CustomerOutput>
{
    private readonly ICustomerRepository _userRepository;

    public CreateCustomer(ICustomerRepository userRepository)
    {
        _userRepository = userRepository;
    }

    public async Task<CustomerOutput> Handle(CreateCustomerInput request, CancellationToken cancellationToken)
    {
        var user = DomainEntity.Customer.Create(request.Name, request.Email);
        await _userRepository.Add(user);
        await _userRepository.Save();
        return CustomerOutput.From(user);
    }
}
