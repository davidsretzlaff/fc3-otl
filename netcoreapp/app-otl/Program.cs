using Microsoft.Data.Sqlite;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using app_otl.Middleware;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
builder.Services.AddControllers();

// Configuração do OpenTelemetry
builder.Services.AddOpenTelemetry()
    .WithTracing(tracerProviderBuilder =>
    {
        tracerProviderBuilder
            // Configura o recurso (informações sobre o serviço)
            .SetResourceBuilder(ResourceBuilder.CreateDefault()
                .AddService("DotNetService"))  // Nome do serviço que aparecerá no Jaeger
            
            // Adiciona instrumentação automática para ASP.NET Core
            // Isso cria spans automaticamente para requisições HTTP
            .AddAspNetCoreInstrumentation()
            
            // Adiciona instrumentação para chamadas HTTP (HttpClient)
            // Isso cria spans para chamadas HTTP que seu serviço faz
            .AddHttpClientInstrumentation()
            
            // Adiciona instrumentação para SQL
            // Isso cria spans para operações no banco de dados
            .AddSqlClientInstrumentation()
            
            // Adiciona o source "MyController" para criar spans manuais
            // Isso permite criar spans personalizados no MyController
            .AddSource("MyController")
            
            // Configura o exportador OTLP para enviar traces ao Jaeger
            .AddOtlpExporter(options =>
            {
                // Endpoint do Jaeger para enviar os traces
                options.Endpoint = new Uri("http://jaeger:4318/v1/traces");
                // Protocolo HTTP com Protobuf
                options.Protocol = OpenTelemetry.Exporter.OtlpExportProtocol.HttpProtobuf;
            });
    });

// Registra o Tracer como um serviço singleton
// Isso permite injetar o Tracer no MyController
builder.Services.AddSingleton<Tracer>(sp => 
    sp.GetRequiredService<TracerProvider>().GetTracer("MyController"));

var app = builder.Build();

// Initialize SQLite database
InitializeDatabase();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseDeveloperExceptionPage();
}

app.UseRouting();

// Adiciona o middleware de tracing
app.UseMiddleware<TracingMiddleware>();

app.UseAuthorization();

app.MapControllers();
app.MapGet("/", () => "Hello World!");

app.Run("http://0.0.0.0:80");

void InitializeDatabase()
{
    using var connection = new SqliteConnection("Data Source=mydatabase.db");
    connection.Open();

    var command = connection.CreateCommand();
    command.CommandText = @"
        DROP TABLE IF EXISTS MyTable;
        CREATE TABLE MyTable (
            Id INTEGER PRIMARY KEY AUTOINCREMENT,
            Name TEXT NOT NULL,
            CreditCard TEXT NOT NULL
        );
    ";
    command.ExecuteNonQuery();

    // Insere dados de teste com cartões de crédito
    command.CommandText = @"
        INSERT INTO MyTable (Name, CreditCard) VALUES 
        ('João Silva', '4532-1234-5678-9012'),
        ('Maria Santos', '5412-8765-4321-9876'),
        ('Pedro Oliveira', '4111-1111-1111-1111');
    ";
    command.ExecuteNonQuery();
}
