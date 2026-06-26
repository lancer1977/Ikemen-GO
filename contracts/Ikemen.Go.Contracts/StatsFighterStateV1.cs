using System.Text.Json.Serialization;

namespace Ikemen.Go.Contracts.V1;

public sealed record StatsFighterStateV1
{
    [JsonPropertyName("name")]
    public string Name { get; init; } = string.Empty;

    [JsonPropertyName("id")]
    public int Id { get; init; }

    [JsonPropertyName("memberNo")]
    public int MemberNo { get; init; }

    [JsonPropertyName("selectNo")]
    public int SelectNo { get; init; }

    [JsonPropertyName("aiLevel")]
    public float AILevel { get; init; }

    [JsonPropertyName("palNo")]
    public int PalNo { get; init; }

    [JsonPropertyName("life")]
    public int Life { get; init; }

    [JsonPropertyName("lifeMax")]
    public int LifeMax { get; init; }

    [JsonPropertyName("power")]
    public int Power { get; init; }

    [JsonPropertyName("powerMax")]
    public int PowerMax { get; init; }

    [JsonPropertyName("winQuote")]
    public int WinQuote { get; init; }

    [JsonPropertyName("win")]
    public bool Win { get; init; }

    [JsonPropertyName("winKO")]
    public bool WinKO { get; init; }

    [JsonPropertyName("winTime")]
    public bool WinTime { get; init; }

    [JsonPropertyName("winPerfect")]
    public bool WinPerfect { get; init; }

    [JsonPropertyName("winClutch")]
    public bool WinClutch { get; init; }

    [JsonPropertyName("winSpecial")]
    public bool WinSpecial { get; init; }

    [JsonPropertyName("winHyper")]
    public bool WinHyper { get; init; }

    [JsonPropertyName("drawGame")]
    public bool DrawGame { get; init; }

    [JsonPropertyName("ko")]
    public bool KO { get; init; }

    [JsonPropertyName("overKO")]
    public bool OverKO { get; init; }
}
