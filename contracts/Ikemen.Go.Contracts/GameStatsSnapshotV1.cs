using System.Text.Json.Serialization;

namespace Ikemen.Go.Contracts.V1;

public interface IGameStatsSnapshotV1
{
    StatsLogV1 StatsLog { get; }
    bool ContinueFlg { get; }
    int PersistRoundCount { get; }
    bool MatchOver { get; }
}

public sealed record GameStatsSnapshotV1 : IGameStatsSnapshotV1
{
    [JsonPropertyName("statsLog")]
    public StatsLogV1 StatsLog { get; init; } = new();

    [JsonPropertyName("continueFlg")]
    public bool ContinueFlg { get; init; }

    [JsonPropertyName("persistRoundCount")]
    public int PersistRoundCount { get; init; }

    [JsonPropertyName("matchOver")]
    public bool MatchOver { get; init; }
}
