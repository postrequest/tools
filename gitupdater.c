#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <dirent.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

void isgit(const char *dirname)
{
    DIR *dirp = opendir(dirname);
    if (dirp != 0) {
        struct dirent *dp;
        while (dp = readdir(dirp), dp != 0) {
            char *path = malloc(strlen(dirname) + 1 + strlen(dp->d_name) + 1);
            if (path) {
                sprintf(path, "%s/%s", dirname, dp->d_name);
                struct stat st;
                if (!strcmp(dp->d_name, ".git")) {
                    if (!lstat(path, &st)) {
                        if (S_ISDIR(st.st_mode)) {
                            pid_t spawn_for_update = fork();
                            if (spawn_for_update == -1) {
                                perror("fork");
                                exit(EXIT_FAILURE);
                            }
                            if (!spawn_for_update) {
                                chdir(dirname);
                                system("git pull");
                                exit(EXIT_SUCCESS);
                            }
                        }
                    }
                }
                free(path);
            }
        }
        closedir(dirp);
    }
}

int main(void)
{
    DIR *dirp;
    if (dirp = opendir("."), dirp != 0) {
        struct dirent *dp;
        while (dp = readdir(dirp), dp != 0) {
            char *path = malloc(2 + strlen(dp->d_name) + 1);
            if (path) {
                sprintf(path, "./%s", dp->d_name);
                if (!strcmp(dp->d_name, ".") && !strcmp(dp->d_name, "..")) {
                    continue;
                }
                struct stat st;
                if (!lstat(path, &st)) {
                    if (S_ISDIR(st.st_mode)) {
                        isgit(path);
                    }
                }
                free(path);
            }
        }
        closedir(dirp);
    }
    while(1) {
        wait(NULL);
        if (errno == ECHILD) {
            break;
        }
    }

    return 0;
}
