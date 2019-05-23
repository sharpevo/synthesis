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
        //static extern GoSlice ReadHumiture();
        static extern ReadHumiture_return ReadHumiture();

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
            ReadHumiture_return humiture = ReadHumiture();
            Console.WriteLine(humiture);
            Console.WriteLine(humiture.r0);
            Console.WriteLine(humiture.r1);
        }
    }
}
