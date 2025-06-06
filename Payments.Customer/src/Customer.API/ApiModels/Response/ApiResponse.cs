namespace app_otl.ApiModels.Response
{
    public class ApiResponse<TData>
    {
        public TData Data { get; private set; }
        public string? Message { get; private set; }
        public string? CorrelationId { get; set; }
        public string? TraceId { get; set; }

        public ApiResponse(TData data)
            => Data = data;

        public ApiResponse(TData data, string? message)
        {
            Data = data;
            Message = message;
        }
    }
}
