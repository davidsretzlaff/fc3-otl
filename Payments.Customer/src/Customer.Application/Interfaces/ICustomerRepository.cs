using Customer.Domain.Domain;

namespace Customer.Application.Interfaces
{
    public interface ICustomerRepository
    {
        Task CreateCustomer(Domain.Domain.Customer user, CancellationToken cancellationToken);
        Task<Domain.Domain.Customer> GetUserById(Guid id, CancellationToken cancellationToken);
        Task Save();
        Task Add(Domain.Domain.Customer user);
    }
}
