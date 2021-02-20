#include "stdio.h"
#include "dirent.h"
#include "stdlib.h"
#include "string.h"
#include <sys/stat.h>
#include <pwd.h>
#include <grp.h>
#include <time.h>

#define MAX_PATH_LEN 1000

struct {
    unsigned int a: 1;
    unsigned int l: 1;
} flags;


int print_path(char* path);
int print_dir(char* path);
int print_file(struct stat* file_stat, char* path);
char** parse_args(int argc, char *argv[], int *num_paths);

int main(int argc, char *argv[]) {

    int num_paths = 0;
    char** paths = parse_args(argc, argv, &num_paths);
//   Had to comment it out because even tho it results in the correct output, the execution ends with a Seg fault.
//   My guess is that when I'm doing paths++ the last time, I'm somehow out of bounds of the allocated memory.
//   That's why I changed it to a for loop, where I can rely on num_paths for the correct number of iterations.
//    char** this_path;
//    while((this_path = paths++)){
//        printf("%s\n", *this_path);
//        if(print_path(*this_path)){
//            exit(EXIT_FAILURE);
//        }
//    }
    for(int i = 0; i < num_paths; i++){
        printf("%s\n", paths[i]);
        if(print_path(paths[i])){
            exit(EXIT_FAILURE);
        }
    }
    exit(EXIT_SUCCESS);
}

char** parse_args(int argc, char *argv[], int *num_paths){
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

    char **paths;
    // Remaining arguments are either dir names or file names
    if(argc > 0) {
        *num_paths = argc;
        paths = (char **)malloc(sizeof(char *) * argc);
        for (int i=0;i<argc;i++)
        {
            paths[i] = (char *)malloc(sizeof(char)*(strlen(*argv+i) + 1));
            strcpy(paths[i], *(argv+i));
        }
        return paths;
    }

    // Use current directory by default
    *num_paths = 1;
    paths = (char **)malloc(sizeof(char *));
    paths[0] = (char *)malloc(sizeof(char) * 2);
    strcpy(paths[0], ".");
    return paths;
}

int print_path(char* path){
    struct stat fileStat;
    stat(path, &fileStat);

    if(S_ISREG(fileStat.st_mode)){
        return print_file(&fileStat, path);
    }

    if(S_ISDIR(fileStat.st_mode)){
        // Assuming that input doesn't have the trailing slash
        return print_dir(strcat(path, "/"));
    }

    printf("only regular files are directories are supported, skipping %s \n", path);
    return 1;
}


int print_dir(char* path) {
    DIR *dir;
    struct dirent *ent;
    if ((dir = opendir(path)) != NULL) {
        struct stat fileStat;
        // print all the files and directories within directory
        while ((ent = readdir(dir)) != NULL) {
            if(!flags.a && (ent->d_name)[0] == '.')
                ;
            else
                if(flags.l) {
                    // Is there a cleaner way to initialize full_path?
                    char* full_path = (char *)malloc(sizeof(char) * (strlen(path) + strlen(ent->d_name) + 1));
                    strcpy(full_path, path);
                    strcat(full_path, ent->d_name);
                    stat(full_path, &fileStat);
                    print_file(&fileStat, ent->d_name);
                } else{
                    printf ("%s\n", ent->d_name);
                }
        }
        closedir (dir);
        printf ("\n");
        return 0;
    } else {
        perror("could not open directory");
        return 1;
    }
}

int print_file(struct stat* fileStat, char* name){
    if(!flags.l) {
        printf ("%s\n", name);
        return 0;
    }
    // File type: "-" for regular or "d" for directory
    printf( (S_ISDIR(fileStat->st_mode)) ? "d" : "-");
    // Read, write, and execution permissions for the file's owner, file's group, everybody else in that order
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

    // File owner's name
    struct passwd *pw = getpwuid(fileStat->st_uid);
    // The name of the group that has file permissions in addition to the file's owner.
    struct group  *gr = getgrgid(fileStat->st_gid);

    if(pw != 0){
        printf("%s ", pw->pw_name);
    }
    if(gr != 0){
        printf("%s ", gr->gr_name);
    }
    // File size in bytes
    printf("%lld ", fileStat->st_size);
    // Last modified
    char buf[80];
    struct tm timestamp = *localtime(&fileStat->st_mtime);
    strftime(buf, 60, "%x %X", &timestamp);
    printf("%s ", buf);
    printf("%s", name);

    printf("\n");
    return 0;
}