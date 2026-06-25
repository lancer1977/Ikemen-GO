using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Ikemen.Go.Contracts.V1;

public sealed record StatsLogV1
{
    [JsonPropertyName("matches")]
    public List<StatsMatchV1> Matches { get; init; } = new();
}
