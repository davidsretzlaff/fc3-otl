using Application.UseCases.User.Common;
using MediatR;

namespace Application.UseCases.User.CreateUser
{
    public interface ICreateUser : IRequestHandler<CreateUserInput, UserOutput>
    {
    }
}
