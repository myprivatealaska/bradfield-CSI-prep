#include "stdio.h"
#include "dirent.h"
#include "stdlib.h"
#include "string.h"
#include <sys/stat.h>
#include "stdbool.h"
#include <pwd.h>
#include <grp.h>

#define MAX_PATH_LEN 1000

struct {
    unsigned int a: 1;
    unsigned int l: 1;
} flags;


int print_dir(const char* path);
void print_full_file_info(struct stat* file_stat, char* path);
void parse_args(char *path, int argc, char *argv[]);

int main(int argc, char *argv[]) {
    // default to current directory
    char path[MAX_PATH_LEN] = ".";

    parse_args(path, argc, argv);

    struct stat fileStat;
    stat(path, &fileStat);

    if(S_ISREG(fileStat.st_mode)){
        print_full_file_info(&fileStat, path);
    }

    if(S_ISDIR(fileStat.st_mode)){
        print_dir(path);
    }

    return EXIT_SUCCESS;
}

void parse_args(char *path, int argc, char *argv[]){
    char c;
    while (--argc > 0 && (*++argv)[0] == '-') {
        while((c = *++argv[0]) != '\0') {
            switch (c) {
                case 'a':
                    flags.a = 1;
                    break;
                case 'l':
                    flags.l = 1;
                    break;
                default:
                    printf("%c flag not implemented, skipping\n", c);
            }
        }
    }

    if(argc > 0) {
        strcpy(path, *argv);
    }
}


int print_dir(const char* path) {
    DIR *dir;
    struct dirent *ent;
    if ((dir = opendir(path)) != NULL) {
        // print all the files and directories within directory
        struct stat fileStat;
        while ((ent = readdir (dir)) != NULL) {
            if(!flags.a && (ent->d_name)[0] == '.')
                ;
            else
                if(flags.l) {
                    stat(ent->d_name, &fileStat);
                    print_full_file_info(&fileStat, ent->d_name);
                } else{
                    printf ("%s\n", ent->d_name);
                }
        }
        closedir (dir);
        return EXIT_SUCCESS;
    } else {
        // could not open directory
        perror ("");
        return 1;
    }
}

void print_full_file_info(struct stat* fileStat, char* path){
    printf( (S_ISDIR(fileStat->st_mode)) ? "d" : "-");
    printf( (fileStat->st_mode & S_IRUSR) ? "r" : "-");
    printf( (fileStat->st_mode & S_IWUSR) ? "w" : "-");
    printf( (fileStat->st_mode & S_IXUSR) ? "x" : "-");
    printf( (fileStat->st_mode & S_IRGRP) ? "r" : "-");
    printf( (fileStat->st_mode & S_IWGRP) ? "w" : "-");
    printf( (fileStat->st_mode & S_IXGRP) ? "x" : "-");
    printf( (fileStat->st_mode & S_IROTH) ? "r" : "-");
    printf( (fileStat->st_mode & S_IWOTH) ? "w" : "-");
    printf( (fileStat->st_mode & S_IXOTH) ? "x" : "-");
    printf(" ");

    // TODO: calculate number of files in the dir

    struct passwd *pw = getpwuid(fileStat->st_uid);
    struct group  *gr = getgrgid(fileStat->st_gid);

    if(pw != 0){
        printf("%s ", pw->pw_name);
    }
    if(gr != 0){
        printf("%s ", gr->gr_name);
    }
    printf("%lld ", fileStat->st_size);
    printf("%ld ", fileStat->st_mtime);
    printf("%s", path);

    printf("\n");
}