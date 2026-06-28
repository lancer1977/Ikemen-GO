using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Ikemen.Go.Contracts.V1;

public sealed record StatsRoundV1
{
    [JsonPropertyName("index")]
    public int Index { get; init; }

    [JsonPropertyName("timer")]
    public int Timer { get; init; }

    [JsonPropertyName("score")]
    public int[] Score { get; init; } = new int[2];

    [JsonPropertyName("fighters")]
    public List<List<StatsFighterStateV1>> Fighters { get; init; } = new()
    {
        new List<StatsFighterStateV1>(),
        new List<StatsFighterStateV1>(),
    };
}
