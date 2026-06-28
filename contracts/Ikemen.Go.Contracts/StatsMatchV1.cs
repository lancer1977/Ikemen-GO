using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Ikemen.Go.Contracts.V1;

public sealed record StatsMatchV1
{
    [JsonPropertyName("matchTime")]
    public int MatchTime { get; init; }

    [JsonPropertyName("roundTime")]
    public int RoundTime { get; init; }

    [JsonPropertyName("winSide")]
    public int WinSide { get; init; }

    [JsonPropertyName("ended")]
    public bool Ended { get; init; }

    [JsonPropertyName("lastRound")]
    public int LastRound { get; init; }

    [JsonPropertyName("draws")]
    public int Draws { get; init; }

    [JsonPropertyName("wins")]
    public int[] Wins { get; init; } = new int[2];

    [JsonPropertyName("teamModes")]
    public int[] TeamModes { get; init; } = new int[2];

    [JsonPropertyName("totalScore")]
    public int[] TotalScore { get; init; } = new int[2];

    [JsonPropertyName("rounds")]
    public List<StatsRoundV1> Rounds { get; init; } = new();
}
