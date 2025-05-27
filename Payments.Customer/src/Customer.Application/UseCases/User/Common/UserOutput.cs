using DomainEntity = Customer.Domain.Domain;

namespace Customer.Application.UseCases.User.Common;

public record UserOutput 
{
    public string Id {get; init;}
    public string Name {get; init;}

    public static UserOutput From(DomainEntity.User user) 
    {
        return new UserOutput()
        {
            Id = user.Id.ToString(),
            Name = user.Name
        };

    }
}
