using System; // for IntPtr
using System.Runtime.InteropServices; // for Mashal
using System.Text; // Encoding

public struct GoString
{
    public string val;
    public int len;
    public static implicit operator GoString(string s)
    {
        return new GoString(){val=s, len=s.Length};
    }
    public static implicit operator string(GoString s){
        //Console.WriteLine("++++");
        //byte[] bytes = new byte[s.n];
        //Marshal.Copy((IntPtr)s.p, bytes, 0, s.n);
        //Console.WriteLine(s.p);
        //Console.WriteLine(s.n);
        //return Encoding.UTF8.GetString(bytes);
        return s.val;

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
