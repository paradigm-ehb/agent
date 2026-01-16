#ifndef PKGMANAGER_H
#define PKGMANAGER_H

#include "arena.h"
#include "base.h"

enum LinuxDistro
{
  DISTRO_UNKNOWN = 0,

  DISTRO_ARCH_LINUX,
  DISTRO_MANJARO,
  DISTRO_ENDEAVOUROS,
  DISTRO_ARCO,
  DISTRO_GARUDA,

  DISTRO_DEBIAN,
  DISTRO_UBUNTU,
  DISTRO_LINUX_MINT,
  DISTRO_POP_OS,

  DISTRO_FEDORA,
  DISTRO_RHEL,
  DISTRO_CENTOS_STREAM,
  DISTRO_ROCKY,
  DISTRO_ALMA,

  DISTRO_GENTOO,
  DISTRO_NIXOS,
  DISTRO_OPENSUSE,

  DISTRO_COUNT

};

struct DistroStack
{
  const char *name;
  enum LinuxDistro linux_distro;
};

internal const struct DistroStack distro_map[] = {
  { "arch", DISTRO_ARCH_LINUX },
  { "manjaro", DISTRO_MANJARO },
  { "endeavouros", DISTRO_ENDEAVOUROS },
  { "arco", DISTRO_ARCO },
  { "garuda", DISTRO_GARUDA },

  { "debian", DISTRO_DEBIAN },
  { "ubuntu", DISTRO_UBUNTU },
  { "mint", DISTRO_LINUX_MINT },
  { "pop", DISTRO_POP_OS },

  { "fedora", DISTRO_FEDORA },
  { "rhel", DISTRO_RHEL },
  { "centos", DISTRO_CENTOS_STREAM },
  { "rocky", DISTRO_ROCKY },
  { "alma", DISTRO_ALMA },

  { "gentoo", DISTRO_GENTOO },
  { "nixos", DISTRO_NIXOS },
  { "opensuse", DISTRO_OPENSUSE },
};

internal enum LinuxDistro
find_lxd_pkgm(mem_arena *arena);

#endif
