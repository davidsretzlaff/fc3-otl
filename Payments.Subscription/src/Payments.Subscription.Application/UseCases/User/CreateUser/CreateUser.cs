using Application.Interfaces;
using Application.UseCases.User.Common;
using DomainEntity = Domain;
using MediatR;

namespace Application.UseCases.User.CreateUser;

public class CreateUser : IRequestHandler<CreateUserInput, UserOutput>
{
    private readonly IUserRepository _userRepository;

    public CreateUser(IUserRepository userRepository)
    {
        _userRepository = userRepository;
    }

    public async Task<UserOutput> Handle(CreateUserInput request, CancellationToken cancellationToken)
    {
        var user = DomainEntity.User.Create(request.Name, request.Login, request.Password);
        await _userRepository.Add(user);
        await _userRepository.Save();
        return UserOutput.From(user);
    }
}
