use std::sync::Arc;

use storage::{storage_sqlite::StorageSqlite, Storage};
use teloxide::prelude::*;

use super::args::ArgsCli;
use super::cmd;
use super::config::Config;
use anyhow::Result;

pub struct App {
    config: Config,
}

impl App {
    pub fn new(args: ArgsCli) -> Self {
        Self {
            config: Config::new(args),
        }
    }

    fn filter_allowed_users(msg: Message, allowed_users_ids: Arc<Vec<u64>>) -> bool {
        allowed_users_ids.contains(&msg.from.unwrap().id.0)
    }
}

impl service::Service for App {
    fn run(&mut self) -> Result<()> {
        pretty_env_logger::init();
        log::info!("Starting MyHealth bot...");

        let bot = Bot::new(self.config.token.clone());

        let handler = dptree::entry().branch(
            Update::filter_message()
                .filter(App::filter_allowed_users)
                .endpoint(cmd::process_command),
        );

        let stg: Arc<Box<dyn Storage>> = Arc::new(Box::new(StorageSqlite::new()?));

        tokio::runtime::Builder::new_multi_thread()
            .enable_all()
            .build()
            .unwrap()
            .block_on(async {
                Dispatcher::builder(bot, handler)
                    .dependencies(dptree::deps![stg, self.config.allowed_user_ids.clone()])
                    .enable_ctrlc_handler()
                    .build()
                    .dispatch()
                    .await;
            });

        Ok(())
    }
}
