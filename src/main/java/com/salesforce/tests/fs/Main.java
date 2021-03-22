package com.salesforce.tests.fs;

import com.salesforce.tests.fs.Path.CmdFilesManager;
import com.salesforce.tests.fs.cmd.CmdInvoker;

import java.util.Arrays;
import java.util.Scanner;

/**
 * The entry point for the Test program
 */
public class Main {

    public static void main(String[] args) {

        Scanner scanner = new Scanner(System.in);
        final CmdInvoker invoker = new CmdInvoker();
        CmdFilesManager.getInstance();
        Boolean isRunning;
        do {

            String[] exec = scanner.nextLine().split(" ");
            isRunning = !exec[0].equals("quit");

            if (isRunning) {
                String[] cmdArgs = null;
                if (exec.length > 1) {
                    cmdArgs = Arrays.copyOfRange(exec, 1, exec.length);
                }
                try {
                    System.out.println(invoker.execute(exec[0]).execute(cmdArgs));
                } catch(RuntimeException e) {
                    System.out.println(e.getMessage());
                }
            }

        } while (isRunning);

    }
}
