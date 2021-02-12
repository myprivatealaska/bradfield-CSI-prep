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


int print_path(char* path);
int print_dir(char* path);
int print_file(struct stat* file_stat, char* path);
void parse_args(char *path, int argc, char *argv[]);

int main(int argc, char *argv[]) {
    // default to current directory
    char path[MAX_PATH_LEN] = ".";

    parse_args(path, argc, argv);

    if(print_path(path)){
        exit(EXIT_FAILURE);
    } else {
        exit(EXIT_SUCCESS);
    }
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

int print_path(char* path){
    struct stat fileStat;
    stat(path, &fileStat);

    if(S_ISREG(fileStat.st_mode)){
        return print_file(&fileStat, path);
    }

    if(S_ISDIR(fileStat.st_mode)){
        return print_dir(path);
    }

    printf("only regular files are directories are supported, skipping %s \n", path);
    return 1;
}


int print_dir(char* path) {
    DIR *dir;
    struct dirent *ent;
    if ((dir = opendir(path)) != NULL) {
        // print all the files and directories within directory
        while ((ent = readdir (dir)) != NULL) {
            if(!flags.a && (ent->d_name)[0] == '.')
                ;
            else
                if(flags.l) {
                    struct stat fileStat;
                    stat(path, &fileStat);
                    print_file(&fileStat, ent->d_name);
                } else{
                    printf ("%s\n", ent->d_name);
                }
        }
        closedir (dir);
        return 0;
    } else {
        perror("could not open directory");
        return 1;
    }
}

int print_file(struct stat* fileStat, char* path){
    if(!flags.l) {
        printf ("%s\n", path);
        return 0;
    }
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
    return 0;
}