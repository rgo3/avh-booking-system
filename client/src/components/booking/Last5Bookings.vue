<template>
  <div>
    <div class="list is-hoverable">
      <div v-for="entry in bookings" :key="entry" class="list-item has-background-white">
        <div class="level">
          <div class="level-left">
            <h6 class="subtitle is-6">{{printDateTime(entry.TimeStamp)}}</h6>
            <br />
            <h6
              class="subtitle is-6"
              v-if="entry.ItemID === 0"
            >{{displayUserName(getUserByID(entry.UserID))}} ~ Payment: {{entry.TotalPrice * -1}}€</h6>
            <h6
              class="subtitle is-6"
              v-else
            >{{displayUserName(getUserByID(entry.UserID))}} ~ {{entry.Amount}}x {{displayItem(getItemByID(entry.ItemID))}}</h6>
          </div>
          <!-- TODO timeout -->
          <button
            v-if="undoableEntries[`entry_${entry.ID}`]"
            class="button is-small"
            @click="undo(entry)"
          >
            <span class="icon is-small">
              <font-awesome-icon icon="redo" />
            </span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
/*
  computed dictionary: entry.ID, bool
    set timeout
  v-if="canUndo(entry.ID)"
    look up bool in dictionary
  doesn't update
*/
export default {
  computed: {
    bookings() {
      var bookings = this.$store.state.last5Bookings;
      this.createUndoableEntry(bookings);
      return bookings;
    }
  },
  data() {
    return {
      undoableEntries: []
    };
  },
  methods: {
    undo(bookEntry) {
      this.$http
        .post("deleteBookEntry", bookEntry)
        .then(() => {
          this.$store.commit("getLast5Bookings");
          var message = "".concat(
            "Deleted book entry from ",
            this.displayUserName(this.getUserByID(bookEntry.UserID))
          );
          this.$responseEventBus.$emit("successMessage", message);
        })
        .catch(() => {
          this.$responseEventBus.$emit(
            "failureMessage",
            "Couldn't undo book entry"
          );
        });
    },
    // isValid(timeStamp) {
    //   // different ts: server vs client
    //   // "2019-08-01T18:27:58Z"
    //   // var current = new Date();
    //   // var ts = new Date(timeStamp);
    //   // console.log({
    //   //   cur: current,
    //   //   ts1: ts.getSeconds(),
    //   //   ts2: ts.setSeconds(ts.getSeconds() + 2)
    //   // });
    //   // if (ts.setSeconds(ts.getSeconds() + 2) < current.getSeconds()) {
    //   //   return false;
    //   // }
    //   // return true;
    //   var nowUnix = new Date().getTime() / 1000;
    //   var tsUnix = new Date(timeStamp).getTime() / 100;
    //   console.log({ now: nowUnix, ts: tsUnix });
    //   if (tsUnix + 5 <= nowUnix) {
    //     return true;
    //   }
    //   return false;
    //   // TODO: refresh
    // }
    createUndoableEntry(bookings) {
      bookings.forEach(entry => {
        var key = `entry_${entry.ID}`;
        this.undoableEntries[key] = false;
      });
      console.log({ a: this.undoableEntries });
      var nowUnix = new Date().getTime() / 1000;
      bookings.forEach(entry => {
        var key = `entry_${entry.ID}`;
        var tsUnix = new Date(entry.TimeStamp).getTime() / 100;
        if (tsUnix + 5 > nowUnix) {
          //TODO 5sec -> 3min
          setTimeout(() => {
            this.undoableEntries[key] = true;
            console.log("time out ", key, this.undoableEntries[key]);
          }, 3000);
          // }, tsUnix + 5 - nowUnix);
        }
      });
      console.log({ b: this.undoableEntries });
      //
      // var nowUnix = new Date().getTime() / 1000;
      // bookings.forEach(entry => {
      //   var tsUnix = new Date(entry.TimeStamp).getTime() / 100;
      //   var key = `entry_${entry.ID}`;
      //   if (tsUnix + 5 > nowUnix) {
      //     //TODO 5sec -> 3min
      //     // undoableEntries.push({ key: true });
      //     setTimeout(() => {
      //       undoableEntries.push({ key: true });
      //     }, tsUnix + 5 - nowUnix);
      //   } else {
      //     undoableEntries.push({ key: false });
      //   }
      //   console.log(undoableEntries);
      // });
      // this.undoableEntries = undoableEntries;
    }
    // canUndo(entryID) {
    //   console.log({
    //     id: entryID,
    //     val: this.undoableEntries[`entry_${entryID}`]
    //   });
    //   return this.undoableEntries[`entry_${entryID}`];
    // }
  }
  // created() {}
};
</script>