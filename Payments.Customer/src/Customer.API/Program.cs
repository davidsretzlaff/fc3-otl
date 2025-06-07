using Customer.Application.Interfaces;
using Customer.Application.UseCases.User.CreateUser;
using Customer.Infra.Data.Repositories; 
using Customer.API.Middleware;
using Microsoft.Data.Sqlite;
using Serilog;
using Serilog.Events;
using Serilog.Formatting;
using System.Text.Json;

var builder = WebApplication.CreateBuilder(args);

// ===== CONFIGURAÇÃO JSON CUSTOMIZADA DO SERILOG =====
Log.Logger = new LoggerConfiguration()
    .MinimumLevel.Information()
    
    // BLOQUEAR COMPLETAMENTE logs do framework/sistema
    .MinimumLevel.Override("Microsoft", LogEventLevel.Fatal)
    .MinimumLevel.Override("Microsoft.AspNetCore", LogEventLevel.Fatal)
    .MinimumLevel.Override("Microsoft.Extensions", LogEventLevel.Fatal)
    .MinimumLevel.Override("Microsoft.Hosting", LogEventLevel.Fatal)
    .MinimumLevel.Override("System", LogEventLevel.Fatal)
    .MinimumLevel.Override("Microsoft.Data", LogEventLevel.Fatal)
    .MinimumLevel.Override("Dapper", LogEventLevel.Fatal)
    
    // PERMITIR apenas logs do nosso domínio
    .MinimumLevel.Override("Customer.API", LogEventLevel.Information)
    .MinimumLevel.Override("Customer.Application", LogEventLevel.Information)
    .MinimumLevel.Override("Customer.Domain", LogEventLevel.Information)
    
    // ENRIQUECER logs com informações de contexto
    .Enrich.FromLogContext()
    .Enrich.WithProperty("service", "customer")
    
    // ESCRITOR DUPLO: Console + Arquivo
    .WriteTo.Console(new CustomJsonFormatter())
    .WriteTo.File(new CustomJsonFormatter(), 
        path: "/app/logs/customer{Date}.log",
        rollingInterval: RollingInterval.Day,
        rollOnFileSizeLimit: true,
        fileSizeLimitBytes: 10485760,
        retainedFileCountLimit: 10)
    
    .CreateLogger();

builder.Host.UseSerilog();

// Add services to the container.
builder.Services.AddControllers();

// Registrar MediatR
builder.Services.AddMediatR(cfg => {
    cfg.RegisterServicesFromAssembly(typeof(CreateCustomerInput).Assembly);
});

// Registrar repositórios e dependências
builder.Services.AddScoped<ICustomerRepository, CustomerRepository>();

var app = builder.Build();

// Initialize SQLite database
InitializeDatabase();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseDeveloperExceptionPage();
}

app.UseRouting();

// MIDDLEWARE LIMPO - apenas correlation ID sem logs redundantes
app.UseCorrelationId();

app.UseAuthorization();

app.MapControllers();

// Health check SILENCIOSO
app.MapGet("/health", () => Results.Ok(new { status = "healthy" }));

app.MapGet("/", () => "Customer Service is running!");

// INICIALIZAÇÃO SILENCIOSA
Log.Information("Customer Service starting up...");
Log.Information("Environment: {Environment}", app.Environment.EnvironmentName);

app.Run("http://0.0.0.0:80");

void InitializeDatabase()
{
    var connectionString = "Data Source=/app/data/customer.db";
    Directory.CreateDirectory("/app/data");
    
    using var connection = new SqliteConnection(connectionString);
    connection.Open();
    
    var createTableCommand = connection.CreateCommand();
    createTableCommand.CommandText = @"
        CREATE TABLE IF NOT EXISTS Customers (
            Id TEXT PRIMARY KEY,
            Name TEXT NOT NULL,
            Email TEXT NOT NULL UNIQUE,
            CreatedAt TEXT NOT NULL,
            UpdatedAt TEXT NOT NULL
        )";
    createTableCommand.ExecuteNonQuery();
    
    // Log APENAS se necessário
    Log.Information("Database initialized successfully");
}

// Formatter customizado para gerar JSON no formato desejado
public class CustomJsonFormatter : ITextFormatter
{
    public void Format(LogEvent logEvent, TextWriter output)
    {
        var logObject = new
        {
            time = logEvent.Timestamp.ToString("yyyy-MM-ddTHH:mm:ssZ"),
            level = logEvent.Level.ToString().ToLower(),
            msg = logEvent.RenderMessage(),
            correlation_id = logEvent.Properties.ContainsKey("correlation_id") 
                ? logEvent.Properties["correlation_id"].ToString().Trim('"') 
                : "",
            service = logEvent.Properties.ContainsKey("service") 
                ? logEvent.Properties["service"].ToString().Trim('"') 
                : "customer"
        };

        var json = JsonSerializer.Serialize(logObject);
        output.WriteLine(json);
    }
}
