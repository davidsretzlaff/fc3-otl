using Application.Interfaces;
using Application.UseCases.User.Common;
using DomainEntity = Domain;

namespace Application.UseCases.User.CreateUser;

public class CreateUser
{
    private readonly IUserRepository _userRepository;

    public CreateUser(IUserRepository userRepository)
    {
        _userRepository = userRepository;
    }

    public async Task<UserOutput> Execute(CreateUserInput input)
    {
        var user = DomainEntity.User.Create(input.Name, input.Login, input.Password);
        await _userRepository.Add(user);
        await _userRepository.Save();
        return UserOutput.From(user);
    }
}
