# JWC

This is a helper program for SWJTU JWC course selection system.

**Disclaimer:** The author is **NOT** responsible to any physical damage,
  mental disorder or broken relationships caused by this program. Use at your
  own risks!

## Usage

- Download corresponding executable for your OS/ARCH
- Edit `config.toml`
- Run `[path_of_the_executable] [path_of_the_config_file]` in the command
  prompt

## Use Cases

### Early Bird

1. Acquire your targets (courseID) before the system is online (usually can be
  acquired during the first selection window)
2. Edit the config file and run the executable few minute before the scheduled
  opening hours
1. If everything is set up correctly, the courses will most likely be chosen
  once the system is online

**NOTICE:** Depending on specific situations (like high parallel number or
  really laggy system), you may chose the same course multiple times at once,
  thus preventing you from selecting other course(s) (also preventing others
  from selecting this course). Therefore, you may need to manually quit
  duplicated courses.

**EXAMPLE CONFIG**

```toml
[client]
delay = 0
parallel = 4
keep = false
```

### Feeling Lucky

1. Run the program in the background until someone decided to quit your
  intended course(s)

**NOTICE:** This is purely luck-based. Definitely not a promising way to choose
  your intended course(s).

**EXAMPLE CONFIG**

```toml
[client]
delay = 1000
parallel = 1
keep = false
```

### Chaotic Evil

1. Find someone who have your intended course(s) chosen and ask him/her to
  provide his/her cookie
2. Prepare 2 instances and set up config file for each account with the same
  targets
3. Run executables simultaneously
4. Ask him/her to quit the course(s)
5. If success, the number of students choosing this course will overflow

**NOTICE:** It may take multiple retries to achieve the goal depending on server
  load. If your intended course is already overflowed, you need at least
  `chosen - planned + 1` accounts to do the trick.

**EXAMPLE CONFIG**

```toml
[client]
delay = 0
parallel = 64
keep = true
```
