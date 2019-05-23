using System; // for IntPtr
using System.Runtime.InteropServices; // for Mashal
using System.Text; // Encoding

public struct GoString
{
    public IntPtr p;
    public int n;
    public static implicit operator GoString(string s)
    {
        Console.WriteLine("----");
        Console.WriteLine(s);
        Console.WriteLine(s.Length);
        return new GoString(){p = Marshal.StringToHGlobalAuto(s), n = s.Length};
    }
    public static implicit operator string(GoString s){
        Console.WriteLine("++++");
        byte[] bytes = new byte[s.n];
        Marshal.Copy((IntPtr)s.p, bytes, 0, s.n);
        Console.WriteLine(s.p);
        Console.WriteLine(s.n);
        return Encoding.UTF8.GetString(bytes);
    }
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
