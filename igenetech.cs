using System;

public struct GoString
{
    public string Value { get; set; }
    public int Length { get; set; }

    public static implicit operator GoString(string s)
    {
        return new GoString() { Value = s, Length = s.Length };
    }

    public static implicit operator string(GoString s) => s.Value;
}
public struct GoSlice
{
    public IntPtr data;
    public int len;
    public int cap;
}

public struct ReadHumiture_return
{
    public double r0;
    public double r1;
}
