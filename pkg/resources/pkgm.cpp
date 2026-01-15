#define DISTRO_FEDORAH
#if defined(DISTRO_ARCH)

#include <libalpm.h>

/**
 * TODO(nasr): pacman impelmentation of package manager
 * */

#elif defined(DISTRO_FEDORAH)

#include <libdnf5/base/base.hpp>
#include <libdnf5/rpm/package_query.hpp>
#include <libdnf5/rpm/transaction_callbacks.hpp>
#include <libdnf5/transaction/transaction_item_action.hpp>
#include <libdnf5/version.hpp>
#include <libdnf5/base/goal.hpp>
#include <libdnf5/repo/package_downloader.hpp>
#include <memory>
#include <vector>

libdnf5::Base *
setup()
{
  libdnf5::Base *base = new libdnf5::Base;
  base->setup();

  auto repo_sack = base->get_repo_sack();

  repo_sack->create_repos_from_system_configuration();
  repo_sack->load_repos();

  return base;
}

std::vector<std::string>
query(libdnf5::Base *base, std::string package_name)
{
  std::vector<std::string> packages;

  libdnf5::rpm::PackageQuery query(*base);
  query.filter_name(package_name);

  for (const auto &pkg : query)
  {
    packages.push_back(pkg.get_nevra());
  }

  return packages;
}

void
run_transaction(libdnf5::Base *base, std::string package_name)
{
  libdnf5::Goal goal(*base);
  goal.add_rpm_install(package_name);

  auto transaction = goal.resolve();

  libdnf5::repo::PackageDownloader downloader(*base);

  for (auto &tspkg : transaction.get_transaction_packages())
  {
    if (libdnf5::transaction::transaction_item_action_is_inbound(tspkg.get_action()))
    {
      downloader.add(tspkg.get_package());
      downloader.download();
    }
  }

  downloader.download();

  transaction.set_callbacks(std::unique_ptr<libdnf5::rpm::TransactionCallbacks>());
  transaction.set_description("installing a package for the first time");

  transaction.run();

  return;
}

#elif defined(DISTRO_DEBIAN)

/**
 * TODO(nasr): apt impelmentation of package manager
 * */

#endif
