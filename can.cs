using System;
using System.Runtime.InteropServices;
using System.Text;

namespace TestCan
{
    class CAN
    {

        //[DllImport("./can.dll", EntryPoint="NewDao")]
        [DllImport("can.dll", EntryPoint="NewDao")]
        static extern int NewDao(byte[] p0, byte[] p1, byte[] p2, byte[] p3, byte[] p4, byte[] p5, byte[] p6, byte[] p7, byte[] p8, byte[] p9);

        //[DllImport("./can.dll", EntryPoint="MoveRelative")]
        [DllImport("can.dll", EntryPoint="MoveRelative")]
        static extern int MoveRelative(int p0, int p1, int p2, int p3);

        [DllImport("can.dll", EntryPoint="ControlSwitcher")]
        static extern int ControlSwitcher(int p0);

        //[DllImport("./can.dll", EntryPoint="ReadHumiture")]
        [DllImport("can.dll", EntryPoint="ReadHumiture")]
        //static extern GoSlice ReadHumiture();
        static extern int ReadHumiture(ref double p0, ref double p1);

        [DllImport("can.dll", EntryPoint="Test")]
        static extern int Test(byte[] p0, ref IntPtr p1);
        //static extern int Test(byte[] p0, ref byte[] p1);

        static void Main()
        {
            IntPtr output= IntPtr.Zero;
            byte[] input = Encoding.UTF8.GetBytes("测试");
            Console.WriteLine("----");
            int len = Test(input, ref output);
            Console.WriteLine("----1"+len);
            //Console.WriteLine("----1.1"+Marshal.PtrToStringUni(output, len));

            byte[] buffer = new byte[len];
            Marshal.Copy(output, buffer, 0, len);
            Console.WriteLine(Encoding.UTF8.GetString(buffer));



            //Console.WriteLine(Marshal.PtrToStringAnsi(output)
            IntPtr output2= IntPtr.Zero;
            byte[] input2 = Encoding.UTF8.GetBytes("0123456");
            len = Test(input2, ref output2);
            Console.WriteLine("----2"+len);

            byte[] buffer2 = new byte[len];
            Marshal.Copy(output2, buffer2, 0, len);
            Console.WriteLine(Encoding.UTF8.GetString(buffer2));

            //byte[] output  = null;
            //byte[] input = Encoding.UTF8.GetBytes("测试");
            //Console.WriteLine("----");
            //int rst = Test(input, ref output);
            //Console.WriteLine("----1"+rst);
            //Console.WriteLine("----1.1"+Encoding.UTF8.GetString(output));
            ////Console.WriteLine(Marshal.PtrToStringAnsi(output)
            //byte[] output2  = null;
            //byte[] input2 = Encoding.UTF8.GetBytes("0123456");
            //rst = Test(input2, ref output2);
            //Console.WriteLine("----2"+rst);
            //Console.WriteLine("----2.1"+Encoding.UTF8.GetString(output2));


//Environment.SetEnvironmentVariable("GODEBUG", "cgocheck=0");
            //int newdao = NewDao(
                    //"4", // devtype
                    //"0", // devindex
                    //"0x00000001", // devid
                    //"0", // canindex
                    //"0x00000000", // acccode
                    //"0xFFFFFFFF", // accmask
                    //"0", // filter
                    //"0x00", // timing0
                    //"0x1c", // timing 1
                    //"0"); // mode
                //Console.WriteLine(newdao);
            ////Console.WriteLine("MoveRelative");
            ////Console.WriteLine(MoveRelative(
                        ////0, 1, 10, 10));
            //Console.WriteLine("ControlSwitcher");
            //Console.WriteLine(ControlSwitcher(12));
            //Console.WriteLine("Readhumiture");
            //double temp = 0; // if float, 32.8 -> 2.93
            //double humi = 0; // if float, 28.1 -> -1.58E23
            //int humiture = ReadHumiture(ref temp, ref humi);
            //Console.WriteLine(humiture);
            //Console.WriteLine(temp);
            //Console.WriteLine(humi);
        }
    }
}
