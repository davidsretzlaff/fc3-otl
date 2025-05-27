namespace Customer.Domain.Domain
{
    public class User
    {
        public Guid Id { get; private set; }
        public string Name { get; private set; }
        public string Login {get; private set;}
        public string Password {get; private set;}
        
        private User() { }
        
        public static User Create(string name, string login, string password)
        {
            var customer = new User
            {
                Id = Guid.NewGuid(),
                Name = name,
                Login = login,
                Password = password,
            };
            return customer;
        }
    }
}
