using System; // for IntPtr
using System.Runtime.InteropServices; // for Mashal
using System.Text; // Encoding

//public struct GoString
//{
    //public string val;
    //public int len;
    //public static implicit operator GoString(string s)
    //{
        //return new GoString(){val=s, len=s.Length};
    //}
    //public static implicit operator string(GoString s){
        ////Console.WriteLine("++++");
        ////byte[] bytes = new byte[s.n];
        ////Marshal.Copy((IntPtr)s.p, bytes, 0, s.n);
        ////Console.WriteLine(s.p);
        ////Console.WriteLine(s.n);
        ////return Encoding.UTF8.GetString(bytes);
        //return s.val;

    //}
//}
//public struct GoString{
    //public IntPtr p;
    //public int n;
    //public static implicit operator GoString(string s){

        //byte[] u8bytes = new byte[s.Length];
        //for (int i = 0; i < s.Length; i++){
            //Console.WriteLine(s[i]);
            //u8bytes[i] = (byte)s[i];
        //}
        //string str = string.Empty;
        //str = Encoding.UTF8.GetString(u8bytes, 0, u8bytes.Length);
        //Console.WriteLine("##");
        //Console.WriteLine(s);
        //Console.WriteLine(str);
        //Console.WriteLine("##");
        //return new GoString(){p=Marshal.StringToHGlobalAuto(str), n=str.Length};
        ////byte[] bytes = Encoding.Default.GetBytes(s);
        ////string utf8 = Encoding.UTF8.GetString(bytes);

        ////return new GoString(){p=Marshal.StringToHGlobalAuto(utf8), n=utf8.Length};
    //}
    //public static implicit operator string(GoString str){
        //byte[] bytes = new byte[str.n];
        //for (int i = 0; i < str.n; i++)
            //bytes[i] = Marshal.ReadByte(str.p, i);
        //string s = Encoding.UTF8.GetString(bytes);
        //return s;
    //}
//}
public struct GoSlice
{
    public IntPtr data;
    public int len;
    public int cap;
}
