const alertmanagerFuncs = [
  "toUpper", "toLower", "title", "trimSpace", "join", "match", "safeHtml", 
  "safeUrl", "urlUnescape", "reReplaceAll", "stringSlice", "date", "tz", 
  "since", "humanizeDuration", "toJson", "list", "append", "dict"
].map(name => ({ label: name, type: "function" }));

const alertmanagerFields = [
  ".Receiver", ".Status", ".Alerts", ".NotificationReason", 
  ".GroupLabels", ".CommonLabels", ".CommonAnnotations", ".ExternalURL"
].map(name => ({ label: name, type: "variable" }));

const alertFields = [
  ".Status", ".Labels", ".Annotations", ".StartsAt", ".EndsAt", 
  ".GeneratorURL", ".Fingerprint"
].map(name => ({ label: name, type: "variable" }));

const alertsMethods = [
  ".Firing", ".Resolved"
].map(name => ({ label: name, type: "method" }));

const kvMethods = [
  ".SortedPairs", ".Names", ".Values"
].map(name => ({ label: name, type: "method" }));

export function createTemplateCompletionSource(alertData) {
  return (context) => {
    // Match the current word, including a leading dot if present
    const word = context.matchBefore(/\.?\w*/);
    
    if (!word || (word.from === word.to && !context.explicit)) {
      return null;
    }

    let options = [];

    // If we have a dot at the start, suggest fields
    if (word.text.startsWith(".")) {
      options = [...alertmanagerFields];
      
      // Try to suggest nested keys if we can parse the path
      // This is a simple implementation: if we match something like .CommonLabels.
      const fullPathMatch = context.matchBefore(/\.[a-zA-Z]+\.[a-zA-Z]*/);
      if (fullPathMatch) {
        const parts = fullPathMatch.text.split("."); // ["", "CommonLabels", ""]
        if (parts.length === 3) {
          const field = parts[1];
          let source = null;
          if (field === "CommonLabels" && alertData?.commonLabels) source = alertData.commonLabels;
          else if (field === "GroupLabels" && alertData?.groupLabels) source = alertData.groupLabels;
          else if (field === "CommonAnnotations" && alertData?.commonAnnotations) source = alertData.commonAnnotations;
          
          if (source) {
            const keys = Object.keys(source).map(key => ({
              label: key,
              type: "property",
              // We don't include a dot because it's already there or we are continuing it
            }));
            
            return {
              from: fullPathMatch.from + field.length + 2, // skip .Field.
              options: keys
            };
          }
        }
      }

      // Also suggest Alert fields and methods if inside a loop? 
      // Hard to know context without a real parser, but let's just add them to the list
      // for now to be helpful, or maybe when user types . after .Alerts
      options = [...options, ...alertFields, ...alertsMethods, ...kvMethods];
    } else {
      // Suggest functions
      options = [...alertmanagerFuncs];
    }

    return {
      from: word.from,
      options: options.filter(o => o.label.startsWith(word.text))
    };
  };
}
