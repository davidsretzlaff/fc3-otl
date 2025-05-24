using Domain;

namespace Application.Interfaces
{
    public interface IUserRepository
    {
        Task CreateUser(User user, CancellationToken cancellationToken);
        Task<User> GetUserById(Guid id, CancellationToken cancellationToken);
        Task Save();
        Task Add(User user);
    }
}
