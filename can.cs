using System;
using System.Runtime.InteropServices;
using System.Text;

namespace TestCan
{
    class Test
    {

        //[DllImport("./can.dll", EntryPoint="NewDao")]
        [DllImport("can.dll", EntryPoint="NewDao")]
        static extern int NewDao(GoString p0, GoString p1, GoString p2, GoString p3, GoString p4, GoString p5, GoString p6, GoString p7, GoString p8, GoString p9);

        //[DllImport("./can.dll", EntryPoint="MoveRelative")]
        [DllImport("can.dll", EntryPoint="MoveRelative")]
        static extern int MoveRelative(int p0, int p1, int p2, int p3);

        [DllImport("can.dll", EntryPoint="ControlSwitcher")]
        static extern int ControlSwitcher(int p0);

        //[DllImport("./can.dll", EntryPoint="ReadHumiture")]
        [DllImport("can.dll", EntryPoint="ReadHumiture")]
        static extern GoSlice ReadHumiture();

        static void Main()
        {
//Environment.SetEnvironmentVariable("GODEBUG", "cgocheck=0");
            int newdao = NewDao(
                    "4", // devtype
                    "0", // devindex
                    "0x00000001", // devid
                    "0", // canindex
                    "0x00000000", // acccode
                    "0xFFFFFFFF", // accmask
                    "0", // filter
                    "0x00", // timing0
                    "0x1c", // timing 1
                    "0"); // mode
                Console.WriteLine(newdao);
            //Console.WriteLine("MoveRelative");
            //Console.WriteLine(MoveRelative(
                        //0, 1, 10, 10));
            Console.WriteLine("ControlSwitcher");
            Console.WriteLine(ControlSwitcher(12));
            Console.WriteLine("Readhumiture");
            GoSlice humiture = ReadHumiture();

//uintptr_t resPtr = phew();
    //uint8_t *res = (uint8_t*)resPtr;

    //for (int i = 0; i < 2; i++){
        //printf("%d\n", res[i]);
    //}

    //printf("Exiting gracefully\n");

            //以下为GoSlice转Array
            //byte[] bytes = new byte[humiture.len];
            //for (int i = 0; i < humiture.len; i++)
                //bytes[i] = Marshal.ReadByte(humiture.data, i);
            ////Byte[] to String
            //string s = Encoding.UTF8.GetString(bytes);
            //Console.WriteLine(s); //输出结果 33 for 23.7, 27.8
            //

            //float[] floats = new float[humiture.len];
            //for (int i = 0; i < humiture.len; i++)
                //floats[i] = Marshal.ReadByte(humiture.data, i);
            ////Byte[] to String
            //string s = Encoding.UTF8.GetString(floats);
            //Console.WriteLine(s); //输出结果 33 for 23.7, 27.8

            double[] array = new double[humiture.len];
            Console.WriteLine(humiture.len);

            Marshal.Copy((IntPtr)humiture.data, array, 0, humiture.len);
            Console.WriteLine(array[0]); //输出结果 33 for 23.7, 27.8
            Console.WriteLine(array[1]); //输出结果 33 for 23.7, 27.8
        }
    }
}
