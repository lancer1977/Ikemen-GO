using System.Text.Json.Serialization;

namespace Ikemen.Go.Contracts.V1;

public interface IGameLiveSnapshotV1 : IGameStatsSnapshotV1
{
    int FrameCounter { get; }
    int MatchTime { get; }
    int CurRoundTime { get; }
    StatsRoundV1 CurrentRound { get; }
}

public sealed record GameLiveSnapshotV1 : IGameLiveSnapshotV1
{
    [JsonPropertyName("statsLog")]
    public StatsLogV1 StatsLog { get; init; } = new();

    [JsonPropertyName("continueFlg")]
    public bool ContinueFlg { get; init; }

    [JsonPropertyName("persistRoundCount")]
    public int PersistRoundCount { get; init; }

    [JsonPropertyName("matchOver")]
    public bool MatchOver { get; init; }

    [JsonPropertyName("frameCounter")]
    public int FrameCounter { get; init; }

    [JsonPropertyName("matchTime")]
    public int MatchTime { get; init; }

    [JsonPropertyName("curRoundTime")]
    public int CurRoundTime { get; init; }

    [JsonPropertyName("currentRound")]
    public StatsRoundV1 CurrentRound { get; init; } = new();
}
