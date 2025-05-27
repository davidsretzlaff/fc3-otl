using DomainEntity = Customer.Domain.Domain;

namespace Customer.Application.UseCases.User.Common;

public record CustomerOutput 
{
    public string Id {get; init;}
    public string Name {get; init;}

    public static CustomerOutput From(DomainEntity.Customer user) 
    {
        return new CustomerOutput()
        {
            Id = user.Id.ToString(),
            Name = user.Name
        };

    }
}
