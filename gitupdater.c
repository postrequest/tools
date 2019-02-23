#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <dirent.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

void update_repo(const char *dirname)
{
    chdir(dirname);
    system("git pull");
}

void isgit(const char *dirname)
{
    DIR *dirp;
    if (dirp = opendir(dirname), dirp != 0) {
	struct dirent *dp;
	while (dp = readdir(dirp), dp != 0) {
	    char *path;
	    if (path = malloc(strlen(dirname) + 1 + strlen(dp->d_name) + 1), path != 0) {
		sprintf(path, "%s/%s", dirname, dp->d_name);

		struct stat st;
		if (strcmp(dp->d_name, ".git") == 0) {
		    if (lstat(path, &st) != -1) {
			    if (S_ISDIR(st.st_mode)) {
                    pid_t spawn_for_update = fork();
                    if (spawn_for_update == -1) {
                        perror("fork");
                        exit(EXIT_FAILURE);
                    }
                    if (spawn_for_update == 0) {
			            update_repo(dirname);
                        exit(EXIT_SUCCESS);
                        wait(NULL);
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
	    char *path;
	    if (path = malloc(strlen(".") + 1 + strlen(dp->d_name) + 1), path != 0) {
		sprintf(path, "%s/%s", ".", dp->d_name);

		struct stat st;
		if (strcmp(dp->d_name, ".") != 0
		    && strcmp(dp->d_name, "..") != 0) {
		    if (lstat(path, &st) != -1) {
                if (S_ISDIR(st.st_mode)) {
                    isgit(path);
                }
		    }
		}

		free(path);
	    }
	}
	closedir(dirp);
    }
    printf("[*] Done\n");

    return 0;
}
