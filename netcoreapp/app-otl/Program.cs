using Microsoft.Data.Sqlite;
using OpenTelemetry.Resources;
using OpenTelemetry.Trace;
using System.Text;
using Serilog;
using Serilog.Events;
using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;

var builder = WebApplication.CreateBuilder(args);

// Configuração do Serilog
Log.Logger = new LoggerConfiguration()
    .MinimumLevel.Information()
    .MinimumLevel.Override("Microsoft", LogEventLevel.Warning)
    .WriteTo.File("logs/app.log", rollingInterval: RollingInterval.Day)
    .CreateLogger();

builder.Host.UseSerilog();

// Add services to the container.
builder.Services.AddControllers();

// Configuração do OpenTelemetry
builder.Services.AddOpenTelemetry()
    .WithTracing(tracerProviderBuilder =>
    {
        tracerProviderBuilder
            .SetResourceBuilder(ResourceBuilder.CreateDefault()
                .AddService("DotNetService")
                .AddTelemetrySdk())
            
            // Configura a instrumentação do ASP.NET Core
            .AddAspNetCoreInstrumentation(options =>
            {
                options.RecordException = true;
                options.Filter = (httpContext) =>
                {
                    // Captura todas as requisições
                    return true;
                };
            })
            
            .AddHttpClientInstrumentation()
            .AddSource("UserController")
            
            .AddOtlpExporter(options =>
            {
                options.Endpoint = new Uri("http://otlcollector:4318/v1/traces");
                options.Protocol = OpenTelemetry.Exporter.OtlpExportProtocol.HttpProtobuf;
            });
    });

builder.Services.AddSingleton<Tracer>(sp => 
    sp.GetRequiredService<TracerProvider>().GetTracer("UserController"));

var app = builder.Build();

// Initialize SQLite database
InitializeDatabase();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseDeveloperExceptionPage();
}

app.UseRouting();

// Middleware para capturar request/response
app.Use(async (context, next) =>
{
    var activity = System.Diagnostics.Activity.Current;
    if (activity != null)
    {
        // Captura o body da requisição
        if (context.Request.Body != null)
        {
            context.Request.EnableBuffering();
            var requestBody = await new StreamReader(context.Request.Body).ReadToEndAsync();
            context.Request.Body.Position = 0;
            activity.SetTag("http.request.body", requestBody);
            activity.SetTag("http.request.path", context.Request.Path);
            activity.SetTag("http.method", context.Request.Method);
        }

        // Captura o body da resposta
        var originalBodyStream = context.Response.Body;
        using var responseBody = new MemoryStream();
        context.Response.Body = responseBody;

        await next();

        responseBody.Position = 0;
        var responseBodyText = await new StreamReader(responseBody).ReadToEndAsync();
        responseBody.Position = 0;
        await responseBody.CopyToAsync(originalBodyStream);
        activity.SetTag("http.response.body", responseBodyText);
        activity.SetTag("http.status_code", context.Response.StatusCode);
    }
    else
    {
        await next();
    }
});

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
